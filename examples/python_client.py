#!/usr/bin/env python3
"""
Example Python client for LLM Gateway using OpenAI SDK.

This demonstrates how to use the gateway with the OpenAI Python client library.
"""

import os
from openai import OpenAI

# Configure OpenAI client to use your gateway
# The virtual key determines which provider and actual API key will be used
client = OpenAI(
    api_key="vk_user1_openai",  # Your virtual key (will route to OpenAI)
    base_url="http://localhost:8080",  # Your gateway URL
)

def test_chat_completion():
    """Test basic chat completion through the gateway."""
    print("Sending chat completion request through gateway...")

    try:
        response = client.chat.completions.create(
            model="gpt-3.5-turbo",
            messages=[
                {"role": "system", "content": "You are a helpful assistant."},
                {"role": "user", "content": "Hello! Can you tell me about Go programming language?"}
            ],
            max_tokens=150
        )

        print("\nResponse received:")
        print(response.choices[0].message.content)

        print(f"\nTokens used: {response.usage.total_tokens}")

    except Exception as e:
        print(f"Error: {e}")

def test_anthropic_via_gateway():
    """
    Test using Anthropic through the gateway.

    Note: This uses a virtual key mapped to Anthropic, but still uses
    the OpenAI SDK format. The gateway handles the provider differences.
    """
    anthropic_client = OpenAI(
        api_key="vk_user2_anthropic",  # Virtual key for Anthropic
        base_url="http://localhost:8080",
    )

    print("\n" + "="*60)
    print("Sending request through gateway to Anthropic...")

    try:
        response = anthropic_client.chat.completions.create(
            model="claude-3-haiku-20240307",
            messages=[
                {"role": "user", "content": "Hello! Can you tell me about Rust programming language?"}
            ],
            max_tokens=150
        )

        print("\nResponse received from Anthropic:")
        print(response.choices[0].message.content)

    except Exception as e:
        print(f"Error: {e}")

def test_streaming():
    """Test streaming responses through the gateway."""
    print("\n" + "="*60)
    print("Testing streaming response...")

    try:
        stream = client.chat.completions.create(
            model="gpt-3.5-turbo",
            messages=[
                {"role": "user", "content": "Count from 1 to 5 slowly."}
            ],
            stream=True
        )

        print("\nStreaming response:")
        for chunk in stream:
            if chunk.choices[0].delta.content is not None:
                print(chunk.choices[0].delta.content, end="", flush=True)
        print()

    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    print("LLM Gateway Client Example")
    print("="*60)

    # Test OpenAI endpoint
    test_chat_completion()

    # Test Anthropic endpoint (commented out by default)
    # Uncomment if you have Anthropic keys configured
    # test_anthropic_via_gateway()

    # Test streaming (commented out by default)
    # test_streaming()

    print("\n" + "="*60)
    print("All tests completed!")
